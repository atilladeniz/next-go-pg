#!/usr/bin/env node
/**
 * Search .docs/ with semantic search (default) or fast fuzzy search
 *
 * Usage:
 *   bun scripts/search-docs.js "how to preload data"     # Semantic (default)
 *   bun scripts/search-docs.js "prefetchQuery" --fast    # Fast fuzzy search
 *   bun scripts/search-docs.js "mutations" --top 3 --llm
 *
 * Options:
 *   --fast, -f       Use fast fuzzy search instead of semantic (for exact keywords)
 *   --top, -n        Number of results (default: 5)
 *   --llm            Output format optimized for LLM context
 *   --verbose, -v    Show debug info
 *
 * Requirements:
 *   bun add -D fuse.js @xenova/transformers
 */

import Fuse from "fuse.js"
import { readdirSync, readFileSync, existsSync, writeFileSync } from "fs"
import { join, relative } from "path"
import { execSync } from "child_process"
import { homedir } from "os"

// Parse arguments
const args = process.argv.slice(2)
let query = ""
let topK = 5
let llmMode = false
let verbose = false
let fast = false // Default is semantic search
let chunkSize = 800
let indexOnly = false

for (let i = 0; i < args.length; i++) {
	const arg = args[i]
	if (arg === "--top" || arg === "-n") {
		topK = parseInt(args[++i], 10)
	} else if (arg === "--llm") {
		llmMode = true
	} else if (arg === "--verbose" || arg === "-v") {
		verbose = true
	} else if (arg === "--fast" || arg === "-f") {
		fast = true
	} else if (arg === "--index") {
		indexOnly = true
	} else if (arg === "--chunk-size" || arg === "-c") {
		chunkSize = parseInt(args[++i], 10)
	} else if (!arg.startsWith("-")) {
		query = arg
	}
}

if (!query && !indexOnly) {
	console.error("Usage: search-docs.js <query> [options]")
	console.error("")
	console.error("Options:")
	console.error("  --fast, -f       Use fast fuzzy search (for exact keywords)")
	console.error("  --top, -n N      Number of results (default: 5)")
	console.error("  --llm            LLM-optimized output format")
	console.error("  --index          Build embeddings index only (no search)")
	console.error("  --verbose, -v    Show debug info")
	console.error("")
	console.error("Examples:")
	console.error('  search-docs.js "how to preload data"        # Semantic (default)')
	console.error('  search-docs.js "prefetchQuery" --fast       # Fast fuzzy search')
	console.error('  search-docs.js "mutations" --top 3 --llm')
	console.error('  search-docs.js --index                      # Pre-build index')
	process.exit(1)
}

// Find .docs directory
function findDocsDir() {
	try {
		const gitRoot = execSync("git rev-parse --show-toplevel", {
			encoding: "utf-8",
		}).trim()
		const docsDir = join(gitRoot, ".docs")
		if (existsSync(docsDir)) return docsDir
	} catch {
		// Not in git repo
	}

	const possiblePaths = [
		join(process.cwd(), ".docs"),
		join(process.cwd(), "..", ".docs"),
	]

	for (const p of possiblePaths) {
		if (existsSync(p)) return p
	}

	throw new Error("Could not find .docs directory")
}

// Recursively get all markdown files
function getMarkdownFiles(dir, files = []) {
	const entries = readdirSync(dir, { withFileTypes: true })

	for (const entry of entries) {
		const fullPath = join(dir, entry.name)
		if (entry.isDirectory()) {
			getMarkdownFiles(fullPath, files)
		} else if (entry.name !== "README.md") {
			// Include .md, .xml, and .txt files
			if (entry.name.endsWith(".md") || entry.name.endsWith(".xml") || entry.name.endsWith(".txt")) {
				files.push(fullPath)
			}
		}
	}

	return files
}

// Split document into chunks with context
function chunkDocument(content, size = 800) {
	const chunks = []
	const sections = content.split(/\n(?=#{1,3}\s)/)

	let currentChunk = ""
	let currentHeader = ""

	for (const section of sections) {
		const headerMatch = section.match(/^(#{1,3}\s+.+?)(?:\n|$)/)
		if (headerMatch) {
			currentHeader = headerMatch[1].trim()
		}

		if (currentChunk.length + section.length <= size) {
			currentChunk += section + "\n"
		} else {
			if (currentChunk.trim()) {
				chunks.push({ text: currentChunk.trim(), header: currentHeader })
			}

			if (section.length > size) {
				const words = section.split(/\s+/)
				let temp = ""
				for (const word of words) {
					if (temp.length + word.length + 1 <= size) {
						temp += word + " "
					} else {
						if (temp.trim()) {
							chunks.push({ text: temp.trim(), header: currentHeader })
						}
						temp = word + " "
					}
				}
				currentChunk = temp
			} else {
				currentChunk = section + "\n"
			}
		}
	}

	if (currentChunk.trim()) {
		chunks.push({ text: currentChunk.trim(), header: currentHeader })
	}

	return chunks
}

// Load and chunk all documents
function loadDocuments(docsDir) {
	const files = getMarkdownFiles(docsDir)
	const passages = []

	for (const file of files) {
		const content = readFileSync(file, "utf-8")
		const relPath = relative(docsDir, file)
		const chunks = chunkDocument(content, chunkSize)

		for (let i = 0; i < chunks.length; i++) {
			passages.push({
				id: `${relPath}#${i}`,
				file: relPath,
				header: chunks[i].header,
				text: chunks[i].text,
				chunkIndex: i,
			})
		}
	}

	return passages
}

// Cosine similarity between two vectors
function cosineSimilarity(a, b) {
	let dotProduct = 0
	let normA = 0
	let normB = 0
	for (let i = 0; i < a.length; i++) {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB))
}

// Get cache file path
function getCachePath(docsDir) {
	return join(docsDir, ".embeddings-cache.json")
}

// Load cached embeddings
function loadCache(cachePath) {
	try {
		if (existsSync(cachePath)) {
			return JSON.parse(readFileSync(cachePath, "utf-8"))
		}
	} catch {
		// Ignore cache errors
	}
	return null
}

// Save embeddings to cache
function saveCache(cachePath, cache) {
	writeFileSync(cachePath, JSON.stringify(cache))
}

// Semantic search with Transformers.js (with caching)
async function semanticSearch(query, passages, topK, verbose) {
	const docsDir = findDocsDir()
	const cachePath = getCachePath(docsDir)

	// Load existing cache
	let cache = loadCache(cachePath)
	const cacheVersion = "v1"

	// Check if cache is valid
	const passageIds = passages.map(p => p.id).sort().join(",")
	const passageHash = passageIds.length.toString() // Simple hash based on count

	if (cache?.version !== cacheVersion || cache?.hash !== passageHash) {
		cache = null
		if (verbose) console.error("📦 Cache invalid, will rebuild...")
	}

	if (verbose) console.error("🧠 Loading embedding model...")

	const { pipeline } = await import("@xenova/transformers")

	// Use a small, fast model
	const embedder = await pipeline("feature-extraction", "Xenova/all-MiniLM-L6-v2", {
		cache_dir: join(homedir(), ".cache", "transformers"),
	})

	// Get query embedding
	if (verbose) console.error("🔍 Embedding query...")
	const queryOutput = await embedder(query, { pooling: "mean", normalize: true })
	const queryEmbedding = Array.from(queryOutput.data)

	// Get or compute passage embeddings
	let passageEmbeddings = cache?.embeddings || {}

	if (!cache) {
		if (verbose) console.error(`🔄 Building embeddings cache for ${passages.length} chunks...`)
		if (verbose) console.error("   (This is a one-time operation, subsequent searches will be fast)")

		for (let i = 0; i < passages.length; i++) {
			const passage = passages[i]
			if (verbose && i % 100 === 0) {
				console.error(`   Processing ${i}/${passages.length}...`)
			}

			const passageOutput = await embedder(passage.text.slice(0, 512), {
				pooling: "mean",
				normalize: true,
			})
			passageEmbeddings[passage.id] = Array.from(passageOutput.data)
		}

		// Save cache
		saveCache(cachePath, {
			version: cacheVersion,
			hash: passageHash,
			embeddings: passageEmbeddings,
		})

		if (verbose) console.error("✅ Cache saved!")
	} else {
		if (verbose) console.error("⚡ Using cached embeddings...")
	}

	// Calculate similarities
	const results = passages.map(passage => ({
		...passage,
		score: cosineSimilarity(queryEmbedding, passageEmbeddings[passage.id]),
	}))

	// Sort by similarity and return top K
	results.sort((a, b) => b.score - a.score)
	return results.slice(0, topK)
}

// Build index only (no search)
async function buildIndex() {
	const docsDir = findDocsDir()
	const cachePath = getCachePath(docsDir)
	const passages = loadDocuments(docsDir)

	console.error(`📁 Docs directory: ${docsDir}`)
	console.error(`📄 Found ${passages.length} chunks to index`)
	console.error("🧠 Loading embedding model...")

	const { pipeline } = await import("@xenova/transformers")

	const embedder = await pipeline("feature-extraction", "Xenova/all-MiniLM-L6-v2", {
		cache_dir: join(homedir(), ".cache", "transformers"),
	})

	console.error(`🔄 Building embeddings for ${passages.length} chunks...`)

	const passageEmbeddings = {}

	for (let i = 0; i < passages.length; i++) {
		const passage = passages[i]
		if (i % 100 === 0) {
			console.error(`   Processing ${i}/${passages.length}...`)
		}

		const passageOutput = await embedder(passage.text.slice(0, 512), {
			pooling: "mean",
			normalize: true,
		})
		passageEmbeddings[passage.id] = Array.from(passageOutput.data)
	}

	const passageIds = passages.map(p => p.id).sort().join(",")
	const passageHash = passageIds.length.toString()

	saveCache(cachePath, {
		version: "v1",
		hash: passageHash,
		embeddings: passageEmbeddings,
	})

	console.error(`✅ Index saved to ${cachePath}`)
	console.error(`📊 Size: ${Math.round(JSON.stringify(passageEmbeddings).length / 1024 / 1024)}MB`)
}

// Fuzzy search with Fuse.js
function fuzzySearch(query, passages, topK, verbose) {
	if (verbose) console.error("🔍 Running fuzzy search...")

	const fuse = new Fuse(passages, {
		keys: [
			{ name: "text", weight: 0.7 },
			{ name: "header", weight: 0.2 },
			{ name: "file", weight: 0.1 },
		],
		includeScore: true,
		threshold: 0.6,
		ignoreLocation: true,
		minMatchCharLength: 2,
		findAllMatches: true,
	})

	const results = fuse.search(query, { limit: topK })

	return results.map((r) => ({
		...r.item,
		score: 1 - (r.score || 0),
	}))
}

// Main search function
async function search(query, topK, verbose, useFast) {
	const docsDir = findDocsDir()
	if (verbose) console.error(`📁 Docs directory: ${docsDir}`)

	const passages = loadDocuments(docsDir)
	if (verbose) console.error(`📄 Loaded ${passages.length} chunks`)

	if (useFast) {
		return fuzzySearch(query, passages, topK, verbose)
	} else {
		return await semanticSearch(query, passages, topK, verbose)
	}
}

// Format for human readability
function formatResult(result, index) {
	const maxLen = 500
	const displayText =
		result.text.length > maxLen
			? result.text.slice(0, maxLen) + "..."
			: result.text

	let output = `━━━ Result ${index + 1} ━━━\n`
	output += `📄 File: ${result.file}\n`
	if (result.header) output += `📌 Section: ${result.header}\n`
	output += `📊 Score: ${result.score.toFixed(4)}\n\n`
	output += displayText + "\n"

	return output
}

// Format for LLM context
function formatForLLM(results, query, useFast) {
	const method = useFast ? "fuzzy" : "semantic"
	let output = `# Search Results for: ${query}\n\n`
	output += `Found ${results.length} relevant sections (${method} search):\n\n`

	for (let i = 0; i < results.length; i++) {
		const r = results[i]
		output += `## [${i + 1}] ${r.file}\n`
		if (r.header) output += `**Section:** ${r.header}\n`
		output += `**Relevance:** ${(r.score * 100).toFixed(0)}%\n\n`
		output += "```\n" + r.text + "\n```\n\n"
	}

	return output
}

// Run
try {
	if (indexOnly) {
		await buildIndex()
		process.exit(0)
	}

	const results = await search(query, topK, verbose, fast)

	if (results.length === 0) {
		console.error("No results found.")
		process.exit(1)
	}

	if (llmMode) {
		console.log(formatForLLM(results, query, fast))
	} else {
		for (let i = 0; i < results.length; i++) {
			console.log(formatResult(results[i], i))
		}
	}
} catch (err) {
	console.error(`Error: ${err.message}`)
	if (verbose) console.error(err.stack)
	process.exit(1)
}

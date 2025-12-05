#!/usr/bin/env node
/**
 * Search .docs/ with fuzzy search using Fuse.js
 *
 * Usage:
 *   node scripts/search-docs.js "your query here"
 *   node scripts/search-docs.js "how to use prefetchQuery" --top 5
 *   node scripts/search-docs.js "mutations" --llm
 *
 * Requirements:
 *   bun add -D fuse.js (or npm install fuse.js)
 */

import Fuse from "fuse.js"
import { readdirSync, readFileSync, existsSync } from "fs"
import { join, relative } from "path"
import { execSync } from "child_process"

// Parse arguments
const args = process.argv.slice(2)
let query = ""
let topK = 5
let llmMode = false
let verbose = false
let chunkSize = 800

for (let i = 0; i < args.length; i++) {
	const arg = args[i]
	if (arg === "--top" || arg === "-n") {
		topK = parseInt(args[++i], 10)
	} else if (arg === "--llm") {
		llmMode = true
	} else if (arg === "--verbose" || arg === "-v") {
		verbose = true
	} else if (arg === "--chunk-size" || arg === "-c") {
		chunkSize = parseInt(args[++i], 10)
	} else if (!arg.startsWith("-")) {
		query = arg
	}
}

if (!query) {
	console.error("Usage: search-docs.js <query> [--top N] [--llm] [--verbose]")
	console.error("")
	console.error("Examples:")
	console.error('  search-docs.js "how to use prefetchQuery"')
	console.error('  search-docs.js "mutations" --top 3')
	console.error('  search-docs.js "authentication" --llm')
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

	// Fallback: check common locations
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
		} else if (
			entry.name.endsWith(".md") &&
			entry.name !== "README.md"
		) {
			files.push(fullPath)
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

			// Handle large sections
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

// Main search
function search(query, topK, verbose) {
	const docsDir = findDocsDir()
	if (verbose) console.error(`📁 Docs directory: ${docsDir}`)

	const passages = loadDocuments(docsDir)
	if (verbose) console.error(`📄 Loaded ${passages.length} chunks`)

	// Configure Fuse.js for best results
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
		score: 1 - (r.score || 0), // Convert to similarity score (higher = better)
	}))
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
function formatForLLM(results, query) {
	let output = `# Search Results for: ${query}\n\n`
	output += `Found ${results.length} relevant sections:\n\n`

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
	const results = search(query, topK, verbose)

	if (results.length === 0) {
		console.error("No results found.")
		process.exit(1)
	}

	if (llmMode) {
		console.log(formatForLLM(results, query))
	} else {
		for (let i = 0; i < results.length; i++) {
			console.log(formatResult(results[i], i))
		}
	}
} catch (err) {
	console.error(`Error: ${err.message}`)
	process.exit(1)
}

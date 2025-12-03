import { spawn } from "bun"
import { mkdir, writeFile } from "fs/promises"

interface DiagnosticError {
	file: string
	line: number
	column: number
	code: string
	message: string
	severity: "error" | "warning"
}

interface TypeCheckReport {
	timestamp: string
	success: boolean
	errorCount: number
	warningCount: number
	errors: DiagnosticError[]
	duration: number
}

async function runTypeCheck(): Promise<void> {
	const startTime = Date.now()

	const proc = spawn(["npx", "tsc", "--noEmit", "--pretty", "false"], {
		cwd: process.cwd(),
		stdout: "pipe",
		stderr: "pipe",
	})

	const stdout = await new Response(proc.stdout).text()
	const stderr = await new Response(proc.stderr).text()
	const output = stdout + stderr

	const errors: DiagnosticError[] = []
	const lines = output.split("\n").filter(Boolean)

	for (const line of lines) {
		// Format: src/app/page.tsx(5,10): error TS2322: Type 'string' is not assignable...
		const match = line.match(/^(.+?)\((\d+),(\d+)\):\s*(error|warning)\s+(TS\d+):\s*(.+)$/)
		if (match) {
			errors.push({
				file: match[1],
				line: parseInt(match[2]),
				column: parseInt(match[3]),
				severity: match[4] as "error" | "warning",
				code: match[5],
				message: match[6],
			})
		}
	}

	const duration = Date.now() - startTime
	const errorCount = errors.filter((e) => e.severity === "error").length
	const warningCount = errors.filter((e) => e.severity === "warning").length

	const report: TypeCheckReport = {
		timestamp: new Date().toISOString(),
		success: errorCount === 0,
		errorCount,
		warningCount,
		errors,
		duration,
	}

	await mkdir(".typecheck", { recursive: true })

	const reportPath = ".typecheck/report.json"
	await writeFile(reportPath, JSON.stringify(report, null, 2))

	// Console output
	const separator = "────────────────────────────────────────"
	console.log("\n📊 TypeCheck Report")
	console.log(separator)

	if (report.success) {
		console.log("✅ No type errors found!")
	} else {
		console.log(`❌ Found ${errorCount} error(s), ${warningCount} warning(s)\n`)

		// Group by file
		const byFile = errors.reduce(
			(acc, err) => {
				if (!acc[err.file]) acc[err.file] = []
				acc[err.file].push(err)
				return acc
			},
			{} as Record<string, DiagnosticError[]>
		)

		for (const [file, fileErrors] of Object.entries(byFile)) {
			console.log(`\n📁 ${file}`)
			for (const err of fileErrors) {
				const icon = err.severity === "error" ? "❌" : "⚠️"
				console.log(`   ${icon} [${err.code}] Line ${err.line}: ${err.message}`)
			}
		}
	}

	console.log(separator)
	console.log(`📄 Report: ${reportPath}\n`)

	process.exit(report.success ? 0 : 1)
}

runTypeCheck()

export default {
	extends: ["@commitlint/config-conventional"],
	rules: {
		// Enforce conventional commit types
		"type-enum": [
			2,
			"always",
			[
				"feat", // New feature
				"fix", // Bug fix
				"docs", // Documentation
				"style", // Formatting (no code change)
				"refactor", // Code restructuring
				"perf", // Performance improvement
				"test", // Adding tests
				"build", // Build system changes
				"ci", // CI configuration
				"chore", // Maintenance
				"revert", // Revert commit
			],
		],
		// Type must be lowercase
		"type-case": [2, "always", "lower-case"],
		// Subject must not be empty
		"subject-empty": [2, "never"],
		// Subject max length
		"subject-max-length": [2, "always", 100],
		// No period at end of subject
		"subject-full-stop": [2, "never", "."],
	},
}

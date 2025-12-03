import reactCompiler from "eslint-plugin-react-compiler"
import tseslint from "typescript-eslint"

export default tseslint.config(
	{
		ignores: ["node_modules/**", ".next/**"],
	},
	{
		files: ["**/*.tsx", "**/*.ts"],
		plugins: {
			"react-compiler": reactCompiler,
		},
		languageOptions: {
			parser: tseslint.parser,
			parserOptions: {
				projectService: true,
			},
		},
		rules: {
			"react-compiler/react-compiler": "error",
		},
	}
)

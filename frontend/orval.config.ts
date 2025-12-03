import { defineConfig } from "orval"

export default defineConfig({
	api: {
		input: {
			target: "../backend/docs/swagger.yaml",
			validation: false,
		},
		output: {
			mode: "tags-split",
			target: "./src/api/endpoints",
			schemas: "./src/api/models",
			client: "react-query",
			httpClient: "fetch",
			baseUrl: "http://localhost:8080/api/v1",
			override: {
				mutator: {
					path: "./src/api/custom-fetch.ts",
					name: "customFetch",
				},
				query: {
					useQuery: true,
					useMutation: true,
				},
			},
		},
	},
})

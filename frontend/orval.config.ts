import { defineConfig } from "orval"

export default defineConfig({
	api: {
		input: "../backend/docs/swagger.yaml",
		output: {
			mode: "tags-split",
			target: "./src/shared/api/endpoints",
			schemas: "./src/shared/api/models",
			client: "react-query",
			httpClient: "fetch",
			baseUrl: "http://localhost:8080/api/v1",
			override: {
				mutator: {
					path: "./src/shared/api/custom-fetch.ts",
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

import path from "node:path"
import type { NextConfig } from "next"

const nextConfig: NextConfig = {
	reactCompiler: true,
	output: "standalone",
	turbopack: {
		root: path.join(__dirname),
	},
	// The @peculiar/asn1-* family relies on runtime class metadata for
	// ASN.1 schema lookup, which turbopack/webpack strip during
	// minification — causing "Cannot get schema for 'iV' target" at build
	// time for routes that import better-auth (the passkey/WebAuthn stack
	// uses asn1-* internally). Externalize them so they load from
	// node_modules at runtime with their decorators intact.
	serverExternalPackages: [
		"@peculiar/asn1-schema",
		"@peculiar/asn1-x509",
		"@peculiar/asn1-ecc",
		"@peculiar/asn1-cms",
		"@peculiar/asn1-pkcs8",
		"@peculiar/asn1-pkcs9",
		"@peculiar/asn1-rsa",
		"@peculiar/asn1-csr",
		"@peculiar/asn1-android",
		"@peculiar/asn1-pfx",
		"@peculiar/asn1-x509-attr",
		"@peculiar/x509",
	],
}

export default nextConfig

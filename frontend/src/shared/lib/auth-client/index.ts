// Public API for shared/lib/auth-client
// Safe to use in "use client" components
export {
	authClient,
	passkey,
	sendVerificationEmail,
	signIn,
	signOut,
	signUp,
	twoFactor,
	useSession,
} from "./auth-client"

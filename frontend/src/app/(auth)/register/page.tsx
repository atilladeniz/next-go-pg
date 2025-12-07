import { redirect } from "next/navigation"

// Mit Magic Link ist keine separate Registrierung nötig
// Neue User werden automatisch beim ersten Login erstellt
export default function RegisterPage() {
	redirect("/login")
}

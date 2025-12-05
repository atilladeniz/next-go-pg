import type { SessionUser } from "../model/types"

type UserInfoProps = {
	user: SessionUser
	showEmail?: boolean
}

export function UserInfo({ user, showEmail = false }: UserInfoProps) {
	return (
		<span className="text-sm text-muted-foreground">
			{user.name || user.email}
			{showEmail && user.name && ` (${user.email})`}
		</span>
	)
}

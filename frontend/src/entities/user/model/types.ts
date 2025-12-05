// User entity types
export type User = {
	id: string
	name: string
	email: string
	image?: string | null
	emailVerified?: boolean
	createdAt?: Date
	updatedAt?: Date
}

export type SessionUser = Pick<User, "id" | "name" | "email">

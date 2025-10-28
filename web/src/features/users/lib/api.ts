import {IUser} from "@/types/auth.types.ts";
import {NoPaginatedResponse} from "@/types/api.types.ts";
import {useAppSession} from "@/utils/session.ts";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
const API_KEY = import.meta.env.VITE_API_KEY || "";

export const usersApi = {
    getUserById: async (userId: string): Promise<NoPaginatedResponse<IUser>> => {
        const session = await useAppSession()
        const response = await fetch(`${API_URL}/users/${userId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${session?.data?.token}`
            },
        })

        if (!response.ok) {
            throw new Error('Failed to login');
        }

        return  await response.json();
    },
}

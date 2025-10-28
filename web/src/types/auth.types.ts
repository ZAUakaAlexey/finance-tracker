export type UserRoleTypes = "admin" | "manager" | "supervisor";

export interface IUser {
  id: number;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  verified_at: string | null;
  email_verified_at: string | null;
  last_login_at: string;
  profile: UserProfile;
  address: UserAddress;
  role: UserRoleTypes;
  roles: UserRole[];
}

interface UserProfile {
  first_name: string | null;
  last_name: string | null;
  full_name: string | null;
  image: string | null;
}

interface UserAddress {
  phone: string;
}

export interface UserRole {
  title: string;
  alias: UserRoleTypes;
  description: string;
}

// interface AuthCredentials {
//   token_type: string;
//   expires_in: number;
//   access_token: string;
//   refresh_token: string;
// }

interface AuthResponseResource {
  user: IUser;
  token: string;
}

export interface AuthResponse {
  data: {
    resource: AuthResponseResource;
  };
  message: string;
  errors: string[];
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface LogoutResponse {
  data: string[];
  message: string;
  errors: string[];
}

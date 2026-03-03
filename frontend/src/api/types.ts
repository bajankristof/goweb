export type AuthWellKnown = {
  providers: AuthProvider[];
};

export type AuthProvider = {
  id: string;
  issuer: string;
  icon?: string;
  name?: string;
};

export type User = {
  userId: string;
  openId: string;
  email: string;
  displayName: string;
};

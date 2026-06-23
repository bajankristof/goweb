export type WellKnownInfo = {
  version: string;
  auth: {
    providers: string[];
  };
};

export type User = {
  id: string;
  email: string;
  displayName: string | null;
  createdAt: string;
  updatedAt: string;
};

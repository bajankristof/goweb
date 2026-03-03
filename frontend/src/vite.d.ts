type ResolvedType<T> = T extends Promise<infer U> ? U : T;

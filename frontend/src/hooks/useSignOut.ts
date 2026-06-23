import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router";

import { signOut } from "../api";

export default function useSignOut() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["auth", "signOut"],
    mutationFn: () => signOut(),
    onSettled: async () => {
      queryClient.clear();
      await navigate("/signin");
    },
  });
}

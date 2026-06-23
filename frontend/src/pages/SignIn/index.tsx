import type { QueryClient } from "@tanstack/react-query";

import { rejectSignedIn } from "../../hooks/guards";
import useWellKnownInfo, {
  wellKnownInfoQuery,
} from "../../hooks/useWellKnownInfo";
import SignInButton from "./components/SignInButton";

import classes from "./index.module.scss";

export function signInLoader(queryClient: QueryClient) {
  const guard = rejectSignedIn(queryClient);
  return async () => {
    await guard();
    await queryClient.ensureQueryData(wellKnownInfoQuery);
    return null;
  };
}

function renderButton(idp: string) {
  const onSignIn = () => {
    window.location.href = `/auth/signin/${encodeURIComponent(idp)}`;
  };

  return <SignInButton key={idp} idp={idp} onSignIn={onSignIn} />;
}

export default function SignIn() {
  const { data: info } = useWellKnownInfo();
  if (!info) {
    return null;
  }

  return (
    <main className={`${classes.SignIn} container-fluid`}>
      <article>
        <h1>Goweb</h1>
        {info.auth.providers.map(renderButton)}
      </article>
    </main>
  );
}

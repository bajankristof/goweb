import type { QueryClient } from "@tanstack/react-query";
import { FaAmazon, FaApple, FaFacebook, FaFingerprint, FaGoogle, FaGithub } from "react-icons/fa6";
import { redirect, useLoaderData } from "react-router";

import { type AuthWellKnown, authWellKnownQuery, currentUserQuery } from "../../api";

import "./index.scss";

type LoaderData = {
  providers?: AuthWellKnown["providers"];
};

export function signInLoader(queryClient: QueryClient) {
  return async (): Promise<LoaderData> => {
    const [wellKnown, user] = await Promise.all([
      queryClient.fetchQuery(authWellKnownQuery),
      queryClient.fetchQuery(currentUserQuery),
    ]);

    if (user) {
      throw redirect("/");
    }

    const { providers } = wellKnown;
    if (providers.length > 1) {
      return { providers };
    }

    window.location.href = providers.length === 1 ? "/auth/signin" : "/auth/authless";
    return {};
  };
}

export default function SignIn() {
  const { providers } = useLoaderData<LoaderData>();
  if (!providers) {
    return <main aria-busy="true"></main>;
  }

  var renderIcon = (id: string) => {
    switch (id) {
      case "amazon":
        return <FaAmazon />;
      case "apple":
        return <FaApple />;
      case "facebook":
        return <FaFacebook />;
      case "google":
        return <FaGoogle />;
      case "github":
        return <FaGithub />;
      default:
        return <FaFingerprint />;
    }
  };

  return (
    <main id="SignIn" className="container">
      <article>
        <h1>Sign In</h1>
        {providers.map((provider) => (
          <button
            onClick={() => {
              window.location.href = `/auth/signin/${provider.id}`;
            }}
          >
            {renderIcon(provider.id)}
            &nbsp;
            {provider.name || provider.id || provider.issuer}
          </button>
        ))}
      </article>
    </main>
  );
}

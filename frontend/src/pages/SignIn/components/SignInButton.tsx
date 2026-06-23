import type { ElementType } from "react";
import { FaFingerprint } from "react-icons/fa6";
import {
  SiApple,
  SiAuthelia,
  SiAuthentik,
  SiFacebook,
  SiGithub,
  SiKeycloak,
  SiOkta,
} from "react-icons/si";
import IconGoogleG from "./IconGoogleG";
import IconPocketId from "./IconPocketId";

export type SignInButtonProps = {
  idp: string;
  onSignIn: () => void;
};

type Config = {
  name: string;
  Icon: ElementType;
};

const CONFIG: Record<string, Config> = {
  apple: { name: "Apple", Icon: SiApple },
  authelia: { name: "Authelia", Icon: SiAuthelia },
  authentik: { name: "Authentik", Icon: SiAuthentik },
  facebook: { name: "Facebook", Icon: SiFacebook },
  github: { name: "GitHub", Icon: SiGithub },
  google: { name: "Google", Icon: IconGoogleG },
  keycloak: { name: "Keycloak", Icon: SiKeycloak },
  okta: { name: "Okta", Icon: SiOkta },
  pocketid: { name: "Pocket ID", Icon: IconPocketId },
};

const FALLBACK: Config = { name: "", Icon: FaFingerprint };

export default function SignInButton({ idp, onSignIn }: SignInButtonProps) {
  const { name, Icon } = CONFIG[idp] ?? { ...FALLBACK, name: idp };
  const isUnknown = !CONFIG[idp];

  return (
    <button
      type="button"
      className={`${idp} ${isUnknown ? "unknown" : ""}`}
      onClick={onSignIn}
    >
      <Icon aria-hidden="true" />
      <span>Sign in with {name}</span>
    </button>
  );
}

import { FaUser } from "react-icons/fa6";
import { NavLink } from "react-router";

import useSignOut from "../../hooks/useSignOut";

import classes from "./index.module.scss";

export default function Header() {
  const signOut = useSignOut();

  return (
    <header className={`${classes.Header} container-fluid`}>
      <nav>
        <ul>
          <li>
            <NavLink to="/">Home</NavLink>
          </li>
        </ul>
        <ul>
          <li>
            <details className="dropdown">
              {/* biome-ignore lint/a11y/useSemanticElements: PicoCSS... */}
              <summary
                role="button"
                aria-label="Account menu"
                className="outline"
              >
                <FaUser />
              </summary>
              <ul dir="rtl">
                <li>
                  <button
                    type="button"
                    role="menuitem"
                    onClick={() => signOut.mutate()}
                    disabled={signOut.isPending}
                  >
                    Sign Out
                  </button>
                </li>
              </ul>
            </details>
          </li>
        </ul>
      </nav>
    </header>
  );
}

import { FaUser } from "react-icons/fa6";
import { NavLink } from "react-router";

import "./index.scss";

export default function Header() {
  return (
    <header id="Header" className="container">
      <nav>
        <ul>
          <li>
            <NavLink to="/">Home</NavLink>
          </li>
        </ul>
        <ul>
          <li>
            <details className="dropdown">
              <summary role="button" className="outline">
                <FaUser />
              </summary>
              <ul dir="rtl">
                <li>
                  <a href="/auth/signout">Sign Out</a>
                </li>
              </ul>
            </details>
          </li>
        </ul>
      </nav>
    </header>
  );
}

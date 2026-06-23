import { useRouteError } from "react-router";

import { APIError, NetworkError } from "../../api";
import classes from "./index.module.scss";

export default function ErrorBoundary() {
  const error = useRouteError();
  console.error(error);

  let title = "Oops!";
  let message = "Sorry, an unexpected error has occurred.";

  if (error instanceof NetworkError) {
    title = "Network Error";
    message = error.message;
  } else if (error instanceof APIError) {
    title = error.statusText;
    message = error.message;
  }

  return (
    <section className={`${classes.ErrorBoundary} container-fluid`}>
      <article>
        <h1>{title}</h1>
        <p>{message}</p>
      </article>
    </section>
  );
}

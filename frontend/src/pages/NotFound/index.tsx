import classes from "./index.module.scss";

export default function NotFound() {
  return (
    <main className={`${classes.NotFound} container-fluid`}>
      <section>
        <h1>Not Found</h1>
        <p>The page you are looking for does not exist.</p>
      </section>
    </main>
  );
}

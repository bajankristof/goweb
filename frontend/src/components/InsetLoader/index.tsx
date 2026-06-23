import classes from "./index.module.scss";

export default function InsetLoader() {
  return (
    <section className={classes.InsetLoader}>
      <div aria-busy="true" />
    </section>
  );
}

import { Link } from "@tanstack/react-router";

const Navbar = () => {
  return (
    <header className="w-full flex flex-row justify-between items-center">
      <section>
        <h1 className="font-semibold text-xl">Auth Flow App</h1>
      </section>
      <section>
        <Link
          to="/auth/login"
          className="w-full bg-blue-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 text-sm"
        >
          Login
        </Link>
      </section>
    </header>
  );
};

export default Navbar;

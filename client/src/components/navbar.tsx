import { Link } from "@tanstack/react-router";
import toast from "react-hot-toast";

const Navbar = () => {
  const isAuthenticated = true;
  const handleLogout = async () => {
    try {
      toast.success("Logout successfull");
    } catch (error) {
      toast.error("An unexpected error occurred. Please try again later.");
    }
  };

  return (
    <header className="w-full flex flex-row justify-between items-center">
      <section>
        <h1 className="font-semibold text-xl">Auth Flow App</h1>
      </section>
      <section className="flex flex-row gap-5 items-center">
        <Link to="/profile">Profile</Link>
        {isAuthenticated ? (
          <button
            className="w-full bg-blue-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 text-sm"
            onClick={handleLogout}
          >
            Logout
          </button>
        ) : (
          <Link
            to="/auth/login"
            className="w-full bg-blue-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 text-sm"
          >
            Login
          </Link>
        )}
      </section>
    </header>
  );
};

export default Navbar;

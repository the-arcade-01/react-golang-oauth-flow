import { Link } from "react-router";
import { useAuth } from "../context/AuthProvider";

const Home = () => {
  const { isAuthenticated, login } = useAuth();

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
        <h2 className="text-3xl font-bold mb-6">Home</h2>
        {isAuthenticated ? (
          <Link
            to="/profile"
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
          >
            Go to Profile
          </Link>
        ) : (
          <button
            onClick={login}
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
          >
            Login with Google
          </button>
        )}
      </div>
    </div>
  );
};

export default Home;

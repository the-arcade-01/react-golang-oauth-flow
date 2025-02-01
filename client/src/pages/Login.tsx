import { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../context/AuthProvider";

const Login = () => {
  const { isAuthenticated, login } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate("/");
    }
  }, [isAuthenticated, navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <h1 className="text-2xl font-bold mb-6 text-center">Login Page</h1>
        <button
          onClick={login}
          className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
        >
          Login with Google
        </button>
      </div>
    </div>
  );
};

export default Login;

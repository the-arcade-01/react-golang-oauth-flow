import { useAuth } from "../context/AuthProvider";
import { useNavigate } from "react-router";

const Profile = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleBackToHome = () => {
    navigate("/");
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <h1 className="text-2xl font-bold mb-4 text-center">
          Welcome, {user?.name}
        </h1>
        <p className="text-lg mb-4 text-center">Email: {user?.email}</p>
        <div className="flex justify-center space-x-4">
          <button
            onClick={handleBackToHome}
            className="bg-gray-500 hover:bg-gray-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
          >
            Back to Home
          </button>
          <button
            onClick={logout}
            className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  );
};

export default Profile;

import { useNavigate } from "react-router";

const Privacy = () => {
  const navigate = useNavigate();

  const handleBackClick = () => {
    navigate("/");
  };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-8 max-w-2xl">
        <button
          onClick={handleBackClick}
          className="mb-4 text-blue-500 underline cursor-pointer"
        >
          Back to Home
        </button>
        <h1 className="text-4xl font-bold mb-4 text-blue-600">
          Privacy Policy
        </h1>
        <p className="mb-4 text-gray-700">
          We value your privacy and want to be transparent about the information
          we access. When you use our service, we only access the following
          profile information:
        </p>
        <ul className="list-disc list-inside mb-4 text-gray-700">
          <li>Google ID</li>
          <li>Name</li>
          <li>Email</li>
          <li>Profile Picture</li>
        </ul>
        <p className="mb-4 text-gray-700">
          We do not access any other information.
        </p>
        <p className="text-gray-700">
          If you have any questions about our privacy practices, please contact
          us at{" "}
          <a
            href="https://arcade.build"
            target="_blank"
            className="text-blue-500 underline"
          >
            arcade.build
          </a>
          .
        </p>
      </div>
    </div>
  );
};

export default Privacy;

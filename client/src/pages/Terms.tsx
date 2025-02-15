import { useNavigate } from "react-router";

const Terms = () => {
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
          Terms and Services
        </h1>
        <p className="mb-4 text-gray-700">
          Welcome to our application. By accessing or using our service, you
          agree to be bound by these terms and conditions.
        </p>
        <h2 className="text-2xl font-semibold mb-2 text-blue-500">
          1. Use of Service
        </h2>
        <p className="mb-4 text-gray-700">
          You agree to use the service only for lawful purposes and in a way
          that does not infringe the rights of others or restrict their use and
          enjoyment of the service.
        </p>
        <h2 className="text-2xl font-semibold mb-2 text-blue-500">
          2. Privacy Policy
        </h2>
        <p className="mb-4 text-gray-700">
          Our privacy policy explains how we collect, use, and protect your
          personal information. By using our service, you consent to our privacy
          policy.
        </p>
        <h2 className="text-2xl font-semibold mb-2 text-blue-500">
          3. Changes to Terms
        </h2>
        <p className="mb-4 text-gray-700">
          We reserve the right to modify these terms at any time. Your continued
          use of the service after any changes indicates your acceptance of the
          new terms.
        </p>
        <h2 className="text-2xl font-semibold mb-2 text-blue-500">
          4. Contact Us
        </h2>
        <p className="mb-4 text-gray-700">
          If you have any questions about these terms, please contact us at{" "}
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

export default Terms;

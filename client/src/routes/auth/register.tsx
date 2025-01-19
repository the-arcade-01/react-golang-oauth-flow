import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { SubmitHandler, useForm } from "react-hook-form";
import toast from "react-hot-toast";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

export const Route = createFileRoute("/auth/register")({
  component: RegisterForm,
});

const registerFormSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

type RegisterFormType = z.infer<typeof registerFormSchema>;

function RegisterForm() {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = useForm<RegisterFormType>({
    resolver: zodResolver(registerFormSchema),
  });

  const onSubmit: SubmitHandler<RegisterFormType> = async (data) => {
    try {
      toast.success("success");
      navigate({
        to: "/auth/login",
      });
    } catch (error) {
      toast.error("An unexpected error occurred. Please try again later.");
    }
  };

  return (
    <div className="flex justify-center items-center min-h-screen">
      <form
        className="lg:w-1/4 sm:w-1/2 p-6 rounded-lg shadow-md bg-white"
        onSubmit={handleSubmit(onSubmit)}
      >
        <p className="text-lg font-semibold text-center">Welcome on board!</p>
        <div className="mb-4">
          <label
            className="block text-gray-700 text-sm font-bold mb-2"
            htmlFor="email"
          >
            Email
          </label>
          <input
            type="text"
            id="email"
            placeholder="eg: john@doe.com"
            {...register("email")}
            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          {errors.email && (
            <span className="text-red-500 text-sm">{errors.email.message}</span>
          )}
        </div>
        <div className="mb-4">
          <label
            className="block text-gray-700 text-sm font-bold mb-2"
            htmlFor="password"
          >
            Password
          </label>
          <input
            type="password"
            id="password"
            {...register("password")}
            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          {errors.password && (
            <span className="text-red-500 text-sm">
              {errors.password.message}
            </span>
          )}
        </div>
        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-blue-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
        >
          Register
        </button>
        <p className="text-gray-500 text-sm text-center pt-4">
          Already have an account?{" "}
          <Link to="/auth/login" className="text-blue-400">
            Login here!
          </Link>
        </p>
      </form>
    </div>
  );
}

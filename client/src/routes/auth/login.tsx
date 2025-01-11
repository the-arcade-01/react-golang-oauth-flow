import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";

export const Route = createFileRoute("/auth/login")({
  component: LoginForm,
});

const loginFormSchema = z.object({
  email: z.string().email(),
  password: z.string().min(4, "Min. length of password 4 chars"),
});

type LoginFormType = z.infer<typeof loginFormSchema>;

function LoginForm() {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormType>({
    resolver: zodResolver(loginFormSchema),
  });

  const onSubmit: SubmitHandler<LoginFormType> = async (data) => {
    await new Promise((resolver) => setTimeout(resolver, 1000));
    console.log(data);
    reset();
  };

  return (
    <div className="flex justify-center items-center min-h-screen">
      <form
        className="w-1/4 p-6 rounded-lg shadow-md bg-white"
        onSubmit={handleSubmit(onSubmit)}
      >
        <p className="text-lg font-semibold text-center">Welcome back!</p>
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
          Login
        </button>
        <p className="text-gray-500 text-sm text-center pt-4">
          Don't have an account?{" "}
          <Link to="/auth/register" className="text-blue-400">
            Register here!
          </Link>
        </p>
      </form>
    </div>
  );
}

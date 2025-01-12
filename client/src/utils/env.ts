import { z } from "zod";

const envSchema = z.object({
  BASE_URL: z.string().nonempty(),
});

export const ENV = envSchema.parse(import.meta.env);

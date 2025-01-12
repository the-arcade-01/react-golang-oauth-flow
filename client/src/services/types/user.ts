import { z } from "zod";

const userSchema = z.object({
  email: z.string().email(),
});

export type User = z.infer<typeof userSchema>;

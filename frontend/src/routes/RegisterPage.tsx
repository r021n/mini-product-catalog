import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Card, CardBody, Input, Button } from "@heroui/react";
import { useAuth } from "../lib/useAuth";
import { ApiError } from "../lib/api";

export function RegisterPage() {
  const { register } = useAuth();
  const nav = useNavigate();

  const [name, setName] = useState("New User");
  const [email, setEmail] = useState("newuser@example.com");
  const [password, setPassword] = useState("password123");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function onSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      await register(name, email, password);
      nav("/");
    } catch (e) {
      if (e instanceof ApiError) setError(e.message);
      else setError("Register failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Card className="max-w-md mx-auto">
      <CardBody className="space-y-4">
        <div className="text-lg font-semibold">Register</div>

        <form onSubmit={onSubmit} className="space-y-3">
          <Input label="Name" value={name} onValueChange={setName} />
          <Input
            label="Email"
            value={email}
            onValueChange={setEmail}
            type="email"
          />
          <Input
            label="Password"
            value={password}
            onValueChange={setPassword}
            type="password"
          />

          {error && <div className="text-sm text-red-600">{error}</div>}

          <Button color="primary" type="submit" isLoading={submitting}>
            Register
          </Button>
        </form>

        <div className="text-sm">
          Sudah punya akun?{" "}
          <Link className="underline" to="/login">
            Login
          </Link>
        </div>
      </CardBody>
    </Card>
  );
}

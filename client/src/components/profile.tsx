import { useAuthStore } from "../store/authStore";

const Profile = () => {
  const token = useAuthStore((state) => state.accessToken);
  return <div>Hi! {token}</div>;
};

export default Profile;

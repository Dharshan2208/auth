import type { User } from '../../types';
import { Badge } from '../common/Badge';

interface ProfileCardProps {
  user: User;
}

export function ProfileCard({ user }: ProfileCardProps) {
  return (
    <div className="bg-white border-2 border-black shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] p-6 sm:p-7">
      <div className="flex items-start justify-between mb-6">
        <h3 className="text-xl font-bold text-black">Profile</h3>
        <Badge variant={user.role === 'admin' ? 'yellow' : 'green'}>
          {user.role}
        </Badge>
      </div>

      <div className="space-y-4">
        <div className="flex items-center gap-4 pb-4 border-b-2 border-gray-100">
          <div className="w-14 h-14 bg-yellow-400 border-2 border-black flex items-center justify-center text-xl font-bold">
            {user.username[0].toUpperCase()}
          </div>
          <div>
            <p className="font-bold text-lg text-black">{user.username}</p>
            <p className="text-sm text-gray-500">User ID: {user.id}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
          <div>
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">
              Email
            </p>
            <p className="text-sm font-medium text-black mt-0.5">
              {user.email}
            </p>
          </div>
          <div>
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">
              Role
            </p>
            <p className="text-sm font-medium text-black mt-0.5 capitalize">
              {user.role}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

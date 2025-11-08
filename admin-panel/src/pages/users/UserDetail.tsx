import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { userService } from '../../services/userService';
import { User, UserWallet, UserActivity } from '../../types';
import Card from '../../components/common/Card';
import Table from '../../components/common/Table';
import { format } from 'date-fns';

const UserDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<User | null>(null);
  const [wallet, setWallet] = useState<UserWallet | null>(null);
  const [activity, setActivity] = useState<UserActivity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      loadUserData();
    }
  }, [id]);

  const loadUserData = async () => {
    if (!id) return;

    setLoading(true);
    try {
      const [userData, walletData, activityData] = await Promise.all([
        userService.getUserById(id),
        userService.getUserWallet(id).catch(() => null),
        userService.getUserActivity(id, 1).catch(() => ({ data: [] })),
      ]);
      setUser(userData);
      setWallet(walletData);
      setActivity(activityData.data || []);
    } catch (error) {
      console.error('Error loading user data:', error);
    } finally {
      setLoading(false);
    }
  };

  const activityColumns = [
    {
      key: 'action',
      label: 'Action',
      render: (activity: UserActivity) => (
        <span className="font-medium text-gray-900 dark:text-white">{activity.action}</span>
      ),
    },
    {
      key: 'details',
      label: 'Details',
      render: (activity: UserActivity) => (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {JSON.stringify(activity.details)}
        </span>
      ),
    },
    {
      key: 'timestamp',
      label: 'Time',
      render: (activity: UserActivity) => format(new Date(activity.timestamp), 'PPp'),
    },
  ];

  if (loading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">User not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{user.username}</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">{user.email}</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card title="User Information">
            <dl className="space-y-4">
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">User ID</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white font-mono">{user.id}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Email</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">{user.email}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Wallet Address</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white font-mono">
                  {user.wallet_address || 'Not connected'}
                </dd>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Joined</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {format(new Date(user.created_at), 'PPpp')}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Last Updated</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {format(new Date(user.updated_at), 'PPpp')}
                  </dd>
                </div>
              </div>
            </dl>
          </Card>

          <Card title="Recent Activity">
            <Table
              columns={activityColumns}
              data={activity}
              keyExtractor={(activity) => activity.id}
              emptyMessage="No recent activity"
            />
          </Card>
        </div>

        <div className="space-y-6">
          <Card title="Wallet Balance">
            {wallet ? (
              <div className="space-y-4">
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Available Balance</dt>
                  <dd className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
                    ${wallet.balance.toLocaleString()}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Locked Balance</dt>
                  <dd className="mt-1 text-xl font-semibold text-gray-700 dark:text-gray-300">
                    ${wallet.locked_balance.toLocaleString()}
                  </dd>
                </div>
                <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Balance</dt>
                  <dd className="mt-1 text-2xl font-bold text-primary-600 dark:text-primary-400">
                    ${(wallet.balance + wallet.locked_balance).toLocaleString()}
                  </dd>
                </div>
              </div>
            ) : (
              <p className="text-gray-500 dark:text-gray-400">No wallet data available</p>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default UserDetail;

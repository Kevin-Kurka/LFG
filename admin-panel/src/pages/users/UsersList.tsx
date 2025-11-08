import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userService } from '../../services/userService';
import { User } from '../../types';
import Table from '../../components/common/Table';
import Input from '../../components/common/Input';
import Pagination from '../../components/common/Pagination';
import Button from '../../components/common/Button';
import { format } from 'date-fns';

const UsersList: React.FC = () => {
  const navigate = useNavigate();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');

  useEffect(() => {
    loadUsers();
  }, [currentPage, search]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const response = await userService.getUsers(currentPage, 20, search);
      setUsers(response.data);
      setTotalPages(response.total_pages);
    } catch (error) {
      console.error('Error loading users:', error);
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    {
      key: 'username',
      label: 'User',
      render: (user: User) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">{user.username}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">{user.email}</div>
        </div>
      ),
    },
    {
      key: 'wallet_address',
      label: 'Wallet',
      render: (user: User) => (
        <span className="text-sm font-mono text-gray-700 dark:text-gray-300">
          {user.wallet_address ? `${user.wallet_address.slice(0, 10)}...` : 'N/A'}
        </span>
      ),
    },
    {
      key: 'created_at',
      label: 'Joined',
      render: (user: User) => format(new Date(user.created_at), 'MMM dd, yyyy'),
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (user: User) => (
        <Button
          size="sm"
          variant="secondary"
          onClick={() => navigate(`/users/${user.id}`)}
        >
          View Details
        </Button>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Users</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Manage platform users
        </p>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="mb-6">
          <Input
            placeholder="Search by username or email..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <Table
          columns={columns}
          data={users}
          keyExtractor={(user) => user.id}
          loading={loading}
          emptyMessage="No users found"
        />

        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={setCurrentPage}
        />
      </div>
    </div>
  );
};

export default UsersList;

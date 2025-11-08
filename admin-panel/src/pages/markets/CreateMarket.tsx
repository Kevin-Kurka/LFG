import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { marketService } from '../../services/marketService';
import { CreateMarketInput } from '../../types';
import Input from '../../components/common/Input';
import Select from '../../components/common/Select';
import Button from '../../components/common/Button';
import Card from '../../components/common/Card';

const CreateMarket: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [outcomes, setOutcomes] = useState<string[]>(['', '']);

  const { register, handleSubmit, formState: { errors } } = useForm<CreateMarketInput>();

  const addOutcome = () => {
    setOutcomes([...outcomes, '']);
  };

  const removeOutcome = (index: number) => {
    if (outcomes.length > 2) {
      setOutcomes(outcomes.filter((_, i) => i !== index));
    }
  };

  const updateOutcome = (index: number, value: string) => {
    const newOutcomes = [...outcomes];
    newOutcomes[index] = value;
    setOutcomes(newOutcomes);
  };

  const onSubmit = async (data: CreateMarketInput) => {
    setError('');
    setLoading(true);

    try {
      const validOutcomes = outcomes.filter((o) => o.trim() !== '');
      if (validOutcomes.length < 2) {
        setError('At least 2 outcomes are required');
        setLoading(false);
        return;
      }

      await marketService.createMarket({
        ...data,
        outcomes: validOutcomes,
      });

      navigate('/markets');
    } catch (err: any) {
      setError(err.message || 'Failed to create market');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Create Market</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Create a new prediction market
        </p>
      </div>

      <Card>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {error && (
            <div className="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
              <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
            </div>
          )}

          <Input
            label="Market Title"
            {...register('title', { required: 'Title is required' })}
            error={errors.title?.message}
            placeholder="Will Bitcoin reach $100k by end of 2025?"
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Description
            </label>
            <textarea
              {...register('description', { required: 'Description is required' })}
              rows={4}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
              placeholder="Provide detailed information about the market..."
            />
            {errors.description && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.description.message}</p>
            )}
          </div>

          <Select
            label="Category"
            {...register('category', { required: 'Category is required' })}
            error={errors.category?.message}
            options={[
              { value: '', label: 'Select a category' },
              { value: 'SPORTS', label: 'Sports' },
              { value: 'POLITICS', label: 'Politics' },
              { value: 'CRYPTO', label: 'Crypto' },
              { value: 'ENTERTAINMENT', label: 'Entertainment' },
              { value: 'OTHER', label: 'Other' },
            ]}
          />

          <Input
            label="Resolution Source"
            {...register('resolution_source', { required: 'Resolution source is required' })}
            error={errors.resolution_source?.message}
            placeholder="e.g., CoinMarketCap, Official Website, etc."
          />

          <Input
            label="End Date"
            type="datetime-local"
            {...register('end_date', { required: 'End date is required' })}
            error={errors.end_date?.message}
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Outcomes (minimum 2)
            </label>
            <div className="space-y-3">
              {outcomes.map((outcome, index) => (
                <div key={index} className="flex space-x-2">
                  <Input
                    value={outcome}
                    onChange={(e) => updateOutcome(index, e.target.value)}
                    placeholder={`Outcome ${index + 1}`}
                  />
                  {outcomes.length > 2 && (
                    <Button
                      type="button"
                      variant="danger"
                      size="sm"
                      onClick={() => removeOutcome(index)}
                    >
                      Remove
                    </Button>
                  )}
                </div>
              ))}
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={addOutcome}
              >
                Add Outcome
              </Button>
            </div>
          </div>

          <div className="flex space-x-4">
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create Market'}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate('/markets')}
            >
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

export default CreateMarket;

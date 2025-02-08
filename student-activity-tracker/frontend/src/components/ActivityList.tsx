import React from 'react';

interface Activity {
  id: string;
  title: string;
  description: string;
  duration: number;
  status: string;
  startTime?: string;
  endTime?: string;
}

interface ActivityListProps {
  activities: Activity[];
  onStart: (id: string) => void;
  onComplete: (id: string) => void;
  onDelete: (id: string) => void;
}

export default function ActivityList({ activities, onStart, onComplete, onDelete }: ActivityListProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'planned':
        return 'bg-yellow-100 text-yellow-800';
      case 'in-progress':
        return 'bg-blue-100 text-blue-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="space-y-4">
      {activities.map((activity) => (
        <div
          key={activity.id}
          className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow"
        >
          <div className="flex justify-between items-start">
            <div>
              <h3 className="text-lg font-semibold">{activity.title}</h3>
              <p className="text-gray-600">{activity.description}</p>
              <div className="mt-2 space-x-2">
                <span className="text-sm text-gray-500">
                  Duration: {activity.duration} minutes
                </span>
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(
                    activity.status
                  )}`}
                >
                  {activity.status}
                </span>
              </div>
              {activity.startTime && (
                <p className="text-sm text-gray-500">
                  Started: {new Date(activity.startTime).toLocaleString()}
                </p>
              )}
              {activity.endTime && (
                <p className="text-sm text-gray-500">
                  Completed: {new Date(activity.endTime).toLocaleString()}
                </p>
              )}
            </div>
            <div className="space-x-2">
              {activity.status === 'planned' && (
                <button
                  onClick={() => onStart(activity.id)}
                  className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                >
                  Start
                </button>
              )}
              {activity.status === 'in-progress' && (
                <button
                  onClick={() => onComplete(activity.id)}
                  className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700"
                >
                  Complete
                </button>
              )}
              <button
                onClick={() => onDelete(activity.id)}
                className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

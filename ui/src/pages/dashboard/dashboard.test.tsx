import { render, screen } from '@testing-library/react';
import { describe, it, expect, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { Dashboard } from './dashboard';
import { useAgents } from '../../store/agentsStore';

describe('Dashboard Component', () => {
  beforeEach(() => {
    useAgents.setState({ agents: {} });
  });

  it('renders the dashboard with no agents', () => {
    render(
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    );

    expect(screen.getByText('Fleet')).toBeInTheDocument();
    expect(screen.getByText('0/0 online')).toBeInTheDocument();
  });

  it('renders with agents', () => {
    useAgents.setState({
      agents: {
        'test-guid': {
          guid: 'test-guid',
          name: 'Test Server',
          online: true,
          last_seen: '12345',
          history: [],
          latest_snapshot: {
            timestamp: 12345,
            data: {
              cpu: [],
              memory: { total: 1024, available: 512, used: 512 },
              disk: { disk_space: { path: '/', total: 100, available: 50, used: 50, usage_percent: 50 }, disk_io_usage: [] },
              network: [],
              agent_ip: '192.168.1.1',
              processes: []
            }
          }
        }
      }
    });

    render(
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    );

    expect(screen.getByText('1/1 online')).toBeInTheDocument();
    expect(screen.getByText('Test Server')).toBeInTheDocument();
  });
});

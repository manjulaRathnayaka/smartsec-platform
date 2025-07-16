import React, { useState, useRef, useEffect } from 'react';
import { useQuery, useMutation } from 'react-query';
import { useAuth } from '../../contexts/AuthContext';
import { mcpAPI } from '../../services/api';
import {
  PaperAirplaneIcon,
  ChatBubbleLeftIcon,
  UserIcon,
  CpuChipIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon
} from '@heroicons/react/24/outline';

const AIAssistant = () => {
  const { user } = useAuth();
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([]);
  const messagesEndRef = useRef(null);

  const { data: queryHistory } = useQuery(
    'mcp-history',
    () => mcpAPI.getQueryHistory(),
    {
      onSuccess: (data) => {
        if (data.data && data.data.length > 0) {
          setMessages(data.data.map(item => ({
            id: item.id,
            type: 'query',
            content: item.query,
            timestamp: new Date(item.timestamp),
            user: true,
            result: item.result
          })));
        }
      }
    }
  );

  const queryMutation = useMutation(
    (query) => mcpAPI.submitQuery(query),
    {
      onSuccess: (data) => {
        const newMessage = {
          id: Date.now(),
          type: 'response',
          content: data.data,
          timestamp: new Date(),
          user: false,
          status: 'success'
        };
        setMessages(prev => [...prev, newMessage]);
      },
      onError: (error) => {
        const errorMessage = {
          id: Date.now(),
          type: 'error',
          content: error.response?.data?.error || 'Something went wrong',
          timestamp: new Date(),
          user: false,
          status: 'error'
        };
        setMessages(prev => [...prev, errorMessage]);
      }
    }
  );

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!message.trim()) return;

    const userMessage = {
      id: Date.now(),
      type: 'query',
      content: message,
      timestamp: new Date(),
      user: true,
      status: 'sending'
    };

    setMessages(prev => [...prev, userMessage]);
    queryMutation.mutate(message);
    setMessage('');
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const MessageBubble = ({ message }) => {
    const isUser = message.user;

    return (
      <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}>
        <div className={`flex ${isUser ? 'flex-row-reverse' : 'flex-row'} items-start space-x-2`}>
          <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
            isUser ? 'bg-blue-500 ml-2' : 'bg-gray-500 mr-2'
          }`}>
            {isUser ? (
              <UserIcon className="w-5 h-5 text-white" />
            ) : (
              <CpuChipIcon className="w-5 h-5 text-white" />
            )}
          </div>

          <div className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
            isUser
              ? 'bg-blue-500 text-white'
              : message.status === 'error'
                ? 'bg-red-100 text-red-800 border border-red-200'
                : 'bg-gray-100 text-gray-800'
          }`}>
            <div className="text-sm">{message.content}</div>

            {message.result && (
              <div className="mt-2 p-2 bg-white bg-opacity-20 rounded text-xs">
                <pre className="whitespace-pre-wrap">{JSON.stringify(message.result, null, 2)}</pre>
              </div>
            )}

            <div className="flex items-center justify-between mt-1">
              <div className="text-xs opacity-70">
                {message.timestamp.toLocaleTimeString()}
              </div>

              {message.status && (
                <div className="flex items-center space-x-1">
                  {message.status === 'sending' && <ClockIcon className="w-3 h-3" />}
                  {message.status === 'success' && <CheckCircleIcon className="w-3 h-3" />}
                  {message.status === 'error' && <XCircleIcon className="w-3 h-3" />}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    );
  };

  const SuggestedQueries = () => {
    const suggestions = [
      "Show all processes that ran with root privileges in the last 24 hours",
      "Which containers were started by non-admin users?",
      "List the most common commands run by marketing team laptops",
      "Show me any suspicious network activity",
      "What are the top 5 processes consuming CPU?",
      "Show all failed login attempts",
      "List containers running on non-standard ports",
      "Show me devices with high memory usage"
    ];

    return (
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <ChatBubbleLeftIcon className="w-5 h-5 mr-2" />
          Suggested Queries
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
          {suggestions.map((suggestion, index) => (
            <button
              key={index}
              onClick={() => setMessage(suggestion)}
              className="text-left text-sm text-blue-600 hover:text-blue-800 hover:bg-blue-50 p-2 rounded-md transition-colors"
            >
              {suggestion}
            </button>
          ))}
        </div>
      </div>
    );
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">AI Security Assistant</h2>
        <p className="text-gray-600 mt-1">
          Ask questions about your security data in natural language
        </p>
      </div>

      {messages.length === 0 && <SuggestedQueries />}

      {/* Chat Messages */}
      <div className="bg-white rounded-lg shadow mb-6">
        <div className="p-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Chat History</h3>
        </div>

        <div className="p-6 h-96 overflow-y-auto">
          {messages.length === 0 ? (
            <div className="text-center text-gray-500 py-12">
              <CpuChipIcon className="w-12 h-12 mx-auto mb-4 text-gray-400" />
              <p>No conversations yet. Ask your first security question!</p>
            </div>
          ) : (
            messages.map((msg) => (
              <MessageBubble key={msg.id} message={msg} />
            ))
          )}

          {queryMutation.isLoading && (
            <div className="flex justify-start mb-4">
              <div className="flex items-start space-x-2">
                <div className="w-8 h-8 rounded-full bg-gray-500 flex items-center justify-center">
                  <CpuChipIcon className="w-5 h-5 text-white" />
                </div>
                <div className="bg-gray-100 px-4 py-2 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <div className="flex space-x-1">
                      <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce"></div>
                      <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                      <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                    </div>
                    <span className="text-sm text-gray-600">Analyzing your query...</span>
                  </div>
                </div>
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input Form */}
      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6">
        <div className="flex space-x-4">
          <div className="flex-1">
            <input
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Ask about your security data... (e.g., Show all processes running as root)"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              disabled={queryMutation.isLoading}
            />
          </div>
          <button
            type="submit"
            disabled={!message.trim() || queryMutation.isLoading}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
          >
            <PaperAirplaneIcon className="w-5 h-5" />
            <span>Send</span>
          </button>
        </div>

        <div className="mt-4 text-sm text-gray-500">
          <p>
            <strong>Examples:</strong> "Show all containers running as root", "Which devices have high CPU usage?",
            "List suspicious network connections", "Show failed authentication attempts"
          </p>
        </div>
      </form>
    </div>
  );
};

export default AIAssistant;

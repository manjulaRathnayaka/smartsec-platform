const express = require('express');
const passport = require('passport');
const OAuth2Strategy = require('passport-oauth2');
const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');
const { body, validationResult } = require('express-validator');
const winston = require('winston');

const router = express.Router();

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.simple(),
  transports: [new winston.transports.Console()]
});

// Mock user database (replace with actual database)
const users = [
  {
    id: 1,
    email: 'admin@smartsec.com',
    password: '$2b$10$8K1p/a0dAoHjDmTGzFWQv.MsQKXvwsQWxs9.2OYuOJn/mGlQXfOSK', // password123
    role: 'admin',
    name: 'Admin User',
    department: 'Security'
  },
  {
    id: 2,
    email: 'user@smartsec.com',
    password: '$2b$10$8K1p/a0dAoHjDmTGzFWQv.MsQKXvwsQWxs9.2OYuOJn/mGlQXfOSK', // password123
    role: 'user',
    name: 'Regular User',
    department: 'IT'
  }
];

// OAuth2 Strategy Configuration
passport.use(new OAuth2Strategy({
  authorizationURL: process.env.OAUTH2_AUTHORIZATION_URL,
  tokenURL: process.env.OAUTH2_TOKEN_URL,
  clientID: process.env.OAUTH2_CLIENT_ID,
  clientSecret: process.env.OAUTH2_CLIENT_SECRET,
  callbackURL: process.env.OAUTH2_CALLBACK_URL
}, async (accessToken, refreshToken, profile, done) => {
  try {
    // In a real application, you would:
    // 1. Fetch user info from the OAuth provider
    // 2. Check if user exists in your database
    // 3. Create or update user record

    // Mock user creation/lookup
    const user = {
      id: profile.id,
      email: profile.email,
      name: profile.name,
      role: 'user', // Default role, can be updated based on your logic
      oauthProvider: 'oauth2'
    };

    return done(null, user);
  } catch (error) {
    return done(error, null);
  }
}));

passport.serializeUser((user, done) => {
  done(null, user.id);
});

passport.deserializeUser(async (id, done) => {
  try {
    const user = users.find(u => u.id === id);
    done(null, user);
  } catch (error) {
    done(error, null);
  }
});

// Generate JWT token
const generateToken = (user) => {
  return jwt.sign(
    {
      id: user.id,
      email: user.email,
      role: user.role,
      name: user.name,
      department: user.department
    },
    process.env.JWT_SECRET || 'your-super-secret-jwt-key',
    { expiresIn: process.env.JWT_EXPIRES_IN || '7d' }
  );
};

// Local login route
router.post('/login', [
  body('email').isEmail().normalizeEmail(),
  body('password').isLength({ min: 6 })
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { email, password } = req.body;

    // Find user
    const user = users.find(u => u.email === email);
    if (!user) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    // Verify password
    const isValidPassword = await bcrypt.compare(password, user.password);
    if (!isValidPassword) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    // Generate JWT token
    const token = generateToken(user);

    logger.info(`User ${user.email} logged in successfully`);

    res.json({
      token,
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        role: user.role,
        department: user.department
      }
    });
  } catch (error) {
    logger.error('Login error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// OAuth2 login route
router.get('/oauth2', passport.authenticate('oauth2'));

// OAuth2 callback route
router.get('/callback',
  passport.authenticate('oauth2', { failureRedirect: '/login' }),
  (req, res) => {
    try {
      const token = generateToken(req.user);

      // Redirect to frontend with token
      res.redirect(`${process.env.CORS_ORIGIN}/auth/callback?token=${token}`);
    } catch (error) {
      logger.error('OAuth callback error:', error);
      res.redirect(`${process.env.CORS_ORIGIN}/login?error=auth_failed`);
    }
  }
);

// Logout route
router.post('/logout', (req, res) => {
  req.logout((err) => {
    if (err) {
      logger.error('Logout error:', err);
      return res.status(500).json({ error: 'Logout failed' });
    }
    res.json({ message: 'Logged out successfully' });
  });
});

// Get current user profile
router.get('/profile', (req, res) => {
  if (!req.user) {
    return res.status(401).json({ error: 'Not authenticated' });
  }

  res.json({
    user: {
      id: req.user.id,
      email: req.user.email,
      name: req.user.name,
      role: req.user.role,
      department: req.user.department
    }
  });
});

module.exports = router;

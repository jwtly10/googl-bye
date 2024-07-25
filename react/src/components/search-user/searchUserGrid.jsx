import {
  CardContent,
  Box,
  Avatar,
  Typography,
  Card,
  Grid,
  Backdrop,
  CircularProgress,
} from '@mui/material';

export default function UsersGrid({ users, onUserSelect, isFindingRepos }) {
  return (
    <Card sx={{ mt: 4, position: 'relative' }}>
      <CardContent>
        <Grid container spacing={3}>
          {users.map((user, index) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={index}>
              <Box
                onClick={() => onUserSelect(user.login)}
                sx={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  p: 2,
                  border: '1px solid #e0e0e0',
                  borderRadius: 2,
                  transition: 'all 0.3s',
                  cursor: 'pointer',
                  '&:hover': {
                    boxShadow: '0 4px 20px 0 rgba(0,0,0,0.12)',
                    transform: 'translateY(-5px)',
                  },
                }}
              >
                <Avatar
                  src={user.avatar_url}
                  alt={user.login}
                  sx={{ width: 80, height: 80, mb: 2 }}
                />
                <Typography variant="h6" gutterBottom>
                  {user.login}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {user.name || 'Name not available'}
                </Typography>
              </Box>
            </Grid>
          ))}
        </Grid>
      </CardContent>
      <Backdrop
        sx={{
          color: '#fff',
          zIndex: (theme) => theme.zIndex.drawer + 1,
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.3)',
          borderRadius: 2,
        }}
        open={isFindingRepos}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
    </Card>
  );
}

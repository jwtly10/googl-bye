import { useEffect } from 'react';
import {
    Box,
    Backdrop,
    TextField,
    Button,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    CircularProgress,
    Typography,
    Card,
    Grid,
} from '@mui/material';

export default function SearchForm({ username, setUsername, handleSearch, isSearching }) {
    useEffect(() => {
        // Load data from localStorage when component mounts
        const searchUser = localStorage.getItem('searchUser');
        if (searchUser) {
            setUsername(JSON.parse(searchUser));
        }
    }, [setUsername]);

    const saveToLocalStorage = (userInput) => {
        localStorage.setItem('searchUser', JSON.stringify(userInput));
    };

    const handleChange = (e) => {
        const { value } = e.target;
        setUsername(value);
        saveToLocalStorage(value);
    };

    return (
        <Box
            position="relative"
            sx={{
                mb: 4,
                filter: isSearching ? 'brightness(0.8)' : 'none',
                transition: 'filter 0.3s ease-in-out',
            }}
        >
            <Card
                sx={{
                    mb: 4,
                    p: 3,
                }}
            >
                <TextField
                    fullWidth
                    margin="normal"
                    placeholder="Github Username"
                    name="username"
                    value={username}
                    onChange={handleChange}
                    variant="outlined"
                    autoComplete="off"
                    InputProps={{
                        style: { textAlign: 'center' },
                    }}
                    inputProps={{
                        style: { textAlign: 'center' },
                    }}
                    sx={{
                        '& .MuiOutlinedInput-notchedOutline': {
                            textAlign: 'center',
                        },
                    }}
                />
                <Grid container justifyContent="center" sx={{ mt: 2 }}>
                    <Grid item xs={12} sm={6} md={4}>
                        <Button
                            onClick={handleSearch}
                            variant="contained"
                            color="primary"
                            fullWidth
                        >
                            Search
                        </Button>
                    </Grid>
                </Grid>
            </Card>
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
                open={isSearching}
            >
                <CircularProgress color="inherit" />
            </Backdrop>
        </Box>
    );
}

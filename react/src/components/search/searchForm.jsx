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

export default function SearchForm({ searchParams, setSearchParams, handleSearch, isSearching }) {
    useEffect(() => {
        // Load data from localStorage when component mounts
        const storedParams = localStorage.getItem('searchParams');
        if (storedParams) {
            setSearchParams(JSON.parse(storedParams));
        }
    }, [setSearchParams]);

    const saveToLocalStorage = (newParams) => {
        localStorage.setItem('searchParams', JSON.stringify(newParams));
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setSearchParams((prevState) => {
            const newState = {
                ...prevState,
                [name]: value,
            };
            saveToLocalStorage(newState);
            return newState;
        });
    };

    const handleOptsChange = (e) => {
        const { name, value } = e.target;
        setSearchParams((prevState) => {
            const newState = {
                ...prevState,
                opts: {
                    ...prevState.opts,
                    [name]: value,
                },
            };
            saveToLocalStorage(newState);
            return newState;
        });
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
                <Typography variant="h4" gutterBottom>
                    Parameters
                </Typography>

                <TextField
                    fullWidth
                    margin="normal"
                    label="Query"
                    name="query"
                    value={searchParams.query}
                    onChange={handleChange}
                />

                <Grid container spacing={1}>
                    <Grid item xs={12} md={6}>
                        <FormControl fullWidth margin="normal">
                            <InputLabel>Sort</InputLabel>
                            <Select
                                name="sort"
                                value={searchParams.opts.sort}
                                onChange={handleOptsChange}
                                label="Sort"
                            >
                                <MenuItem value="stars">Stars</MenuItem>
                                <MenuItem value="forks">Forks</MenuItem>
                                <MenuItem value="updated">Updated</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <FormControl fullWidth margin="normal">
                            <InputLabel>Order</InputLabel>
                            <Select
                                name="order"
                                value={searchParams.opts.order}
                                onChange={handleOptsChange}
                                label="Order"
                            >
                                <MenuItem value="asc">Ascending</MenuItem>
                                <MenuItem value="desc">Descending</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>

                <Grid container spacing={2}>
                    <Grid item xs={12} md={6}>
                        <TextField
                            fullWidth
                            margin="normal"
                            type="number"
                            label="Start Page"
                            name="startPage"
                            value={searchParams.startPage}
                            onChange={handleChange}
                        />
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            fullWidth
                            margin="normal"
                            type="number"
                            label="Pages to Process*"
                            name="pagesToProcess"
                            value={searchParams.pagesToProcess}
                            onChange={handleChange}
                        />
                    </Grid>
                </Grid>

                <Grid container justifyContent="space-between" sx={{ mt: 2 }}>
                    <Grid item>
                        <Typography variant="caption">
                            *Each page from GitHub can have a maximum of 100 results, with a total limit of 2000
                            results.
                        </Typography>
                    </Grid>
                    <Grid item>
                        <Button onClick={handleSearch} variant="contained" color="primary">
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

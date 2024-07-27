import React from 'react';
import { Box, Container, Typography, Link } from '@mui/material';
import Iconify from '../iconify/iconify';

function Footer() {
    return (
        <Box
            component="footer"
            sx={{
                py: 3,
                px: 2,
                mt: 'auto',
                backgroundColor: (theme) =>
                    theme.palette.mode === 'light'
                        ? theme.palette.grey[200]
                        : theme.palette.grey[800],
            }}
        >
            <Container maxWidth="lg">
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: { xs: 'column', sm: 'row' },
                        justifyContent: 'space-between',
                        alignItems: 'center',
                    }}
                >
                    <Typography variant="body2" color="text.secondary">
                        © {new Date().getFullYear()}{' '}
                        <Link href="https://github.com/jwtly10" target="_blank">
                            jwtly10
                        </Link>
                        . All rights reserved.
                    </Typography>
                    <Box>
                        <Link href="#" color="inherit" sx={{ mx: 1 }}>
                            Privacy Policy
                        </Link>
                        <Link href="#" color="inherit" sx={{ mx: 1 }}>
                            Terms of Service
                        </Link>
                    </Box>
                    <Box>
                        <Link
                            href="https://github.com/jwtly10/googl-bye"
                            color="inherit"
                            target="_blank"
                            sx={{ mx: 1 }}
                        >
                            <Iconify icon="mdi:github" width={24} height={24} />
                        </Link>
                    </Box>
                </Box>
            </Container>
        </Box>
    );
}

export default Footer;

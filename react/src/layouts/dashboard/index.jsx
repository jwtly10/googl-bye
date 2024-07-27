import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Box } from '@mui/material';
import Header from './header';
import Footer from 'src/components/footer/footer';
import Main from './Main'; // Assuming you have this component

export default function DashboardLayout({ children }) {
    const [openNav, setOpenNav] = useState(false);
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                minHeight: '100vh', // This ensures the layout takes at least the full viewport height
            }}
        >
            <Header onOpenNav={() => setOpenNav(true)} />
            <Box
                sx={{
                    flex: 1, // This allows the content to grow and push the footer down
                    display: 'flex',
                    flexDirection: { xs: 'column', lg: 'row' },
                }}
            >
                <Main>{children}</Main>
            </Box>
            <Footer />
        </Box>
    );
}

DashboardLayout.propTypes = {
    children: PropTypes.node,
};

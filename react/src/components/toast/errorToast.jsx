import { Snackbar, Alert } from '@mui/material';
import { useState, useEffect } from 'react';

export default function ErrorToast({ message, open, onClose }) {
    const [isOpen, setIsOpen] = useState(open);

    useEffect(() => {
        setIsOpen(open);
    }, [open]);

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        setIsOpen(false);
        if (onClose) {
            onClose();
        }
    };

    return (
        <Snackbar
            open={isOpen}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            sx={{ zIndex: (theme) => theme.zIndex.drawer + 2 }}
        >
            <Alert
                onClose={handleClose}
                severity="error"
                variant="filled"
                sx={{
                    width: '100%',
                    bgcolor: 'error.dark',
                    color: 'white',
                    '& .MuiAlert-icon': {
                        color: 'white',
                    },
                }}
            >
                {message}
            </Alert>
        </Snackbar>
    );
}

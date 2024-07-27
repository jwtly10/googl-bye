import { useState } from 'react';
import PropTypes from 'prop-types';

import {
    Popover,
    TableRow,
    MenuItem,
    TableCell,
    Typography,
    IconButton,
    Alert,
    Collapse,
    Box,
    Paper,
    TableContainer,
    Table,
    Link,
    TableHead,
    Tooltip,
    TableBody,
} from '@mui/material';
import { styled } from '@mui/material/styles';

import Label from 'src/components/label';
import Iconify from 'src/components/iconify';

export default function IssueTableRow({
    selected,
    name,
    author,
    language,
    stars,
    ghUrl,
    forks,
    state,
    errorMsg,
    issues,
    handleCreateIssue,
}) {
    const [open, setOpen] = useState(null);
    const [expandOpen, setExpandOpen] = useState(false);

    const handleOpenMenu = (event) => {
        setOpen(event.currentTarget);
    };

    const handleCloseMenu = () => {
        setOpen(null);
    };

    const handleExpandToggle = () => {
        setExpandOpen(!expandOpen);
    };

    const StyledTableCell = styled(TableCell)({
        maxWidth: '200px',
        wordWrap: 'break-word',
        overflowWrap: 'break-word',
    });

    const handleCreatingIssue = () => {
        setOpen(null);
        handleCreateIssue();
    };

    return (
        <>
            <TableRow hover tabIndex={-1} role="checkbox" selected={selected}>
                <TableCell
                    padding="checkbox"
                    sx={{
                        width: 48,
                        paddingLeft: 1,
                        paddingRight: 1,
                    }}
                >
                    <IconButton
                        size="small"
                        onClick={(event) => {
                            event.stopPropagation();
                            handleExpandToggle();
                        }}
                    >
                        <Iconify
                            icon={
                                expandOpen
                                    ? 'eva:arrow-ios-upward-fill'
                                    : 'eva:arrow-ios-downward-fill'
                            }
                        />
                    </IconButton>
                </TableCell>
                <TableCell>
                    <Link
                        href={ghUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        underline="hover"
                        color="inherit"
                    >
                        <Typography variant="subtitle2" noWrap>
                            {name}
                        </Typography>
                    </Link>
                </TableCell>
                <TableCell>{author}</TableCell>
                <TableCell>{language}</TableCell>
                <TableCell align="center">{stars}</TableCell>
                <TableCell align="center">{forks}</TableCell>
                <TableCell>
                    <Label
                        color={
                            (state === 'PENDING' && 'warning') ||
                            (state === 'PROCESSING' && 'info') ||
                            (state === 'COMPLETED' && 'success') ||
                            (state === 'ERROR' && 'error') ||
                            'default'
                        }
                    >
                        {state}
                    </Label>
                </TableCell>
                <TableCell align="center">{issues.length}</TableCell>
                <TableCell align="right">
                    <IconButton onClick={handleOpenMenu}>
                        <Iconify icon="eva:more-vertical-fill" />
                    </IconButton>
                </TableCell>
            </TableRow>

            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={9}>
                    <Collapse in={expandOpen} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 2 }}>
                            {errorMsg && (
                                <Alert severity="error" variant="outlined" sx={{ mb: 2 }}>
                                    {errorMsg}
                                </Alert>
                            )}
                            <Typography variant="h6" gutterBottom component="div">
                                Goo.gl Links Found
                            </Typography>
                            {issues.length > 0 ? (
                                <TableContainer component={Paper} elevation={0} variant="outlined">
                                    <Table size="small" aria-label="links">
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Short URL</TableCell>
                                                <TableCell>Expanded URL</TableCell>
                                                <TableCell>File</TableCell>
                                                <TableCell>Line Number</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {issues.map((link) => (
                                                <TableRow key={link.id}>
                                                    <TableCell>
                                                        <Link
                                                            href={link.url}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                        >
                                                            {link.url}
                                                        </Link>
                                                    </TableCell>
                                                    <StyledTableCell>
                                                        {link.expandedUrl.startsWith('ERROR:') ? (
                                                            <Typography color="error">
                                                                {link.expandedUrl}
                                                            </Typography>
                                                        ) : (
                                                            <Link
                                                                href={link.expandedUrl}
                                                                target="_blank"
                                                                rel="noopener noreferrer"
                                                            >
                                                                {link.expandedUrl}
                                                            </Link>
                                                        )}
                                                    </StyledTableCell>
                                                    <StyledTableCell>
                                                        <Link
                                                            href={link.githubUrl}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                        >
                                                            {link.file}
                                                        </Link>
                                                    </StyledTableCell>
                                                    <TableCell>{link.lineNumber}</TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            ) : (
                                <Typography variant="body2" color="text.secondary">
                                    No goo.gl links found in this repository.
                                </Typography>
                            )}
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>

            <Popover
                open={!!open}
                anchorEl={open}
                onClose={handleCloseMenu}
                anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
                transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                PaperProps={{
                    sx: { width: 170 },
                }}
            >
                <MenuItem onClick={handleCreatingIssue} sx={{ mr: 0 }}>
                    <Iconify icon="mdi:github" width={20} height={20} sx={{ mr: 1 }} />
                    Create Issue
                </MenuItem>
                <MenuItem onClick={handleCloseMenu} sx={{ mr: 1 }}>
                    <Iconify icon="mdi:github" width={20} height={20} sx={{ mr: 1 }} />
                    Create PR
                </MenuItem>
                <MenuItem onClick={handleCloseMenu} sx={{ color: 'error.main' }}>
                    <Iconify icon="eva:trash-2-outline" width={20} height={20} sx={{ mr: 1 }} />
                    Delete
                </MenuItem>
            </Popover>
        </>
    );
}

IssueTableRow.propTypes = {
    author: PropTypes.string,
    language: PropTypes.string,
    name: PropTypes.string,
    ghUrl: PropTypes.string,
    stars: PropTypes.number,
    forks: PropTypes.number,
    selected: PropTypes.bool,
    state: PropTypes.string,
    issues: PropTypes.array,
};

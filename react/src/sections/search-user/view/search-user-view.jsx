import { useEffect, useState } from 'react';

import {
    Stack,
    Table,
    Container,
    TableBody,
    TableRow,
    TableCell,
    TableContainer,
    TablePagination,
    Typography,
    Card,
} from '@mui/material';

import Scrollbar from 'src/components/scrollbar';

import SearchUserForm from 'src/components/search-user/searchUserForm';
import ErrorToast from 'src/components/toast/errorToast';
import SuccessToast from 'src/components/toast/successToast';

import TableNoData from '../table-no-data';

import IssueTableRow from '../issues-table-row';
import IssueTableHead from '../issues-table-head';
import IssueTableToolbar from '../issues-table-toolbar';

import TableEmptyRows from '../table-empty-rows';
import { emptyRows, applyFilter, getComparator } from '../utils';

import {
    searchGithubUsers,
    searchGithubUsersRepo,
    saveRepos,
    searchRepoLinksForUser,
    createIssue,
} from 'src/api/client';
import UserGrid from 'src/components/search-user/searchUserGrid';

// ----------------------------------------------------------------------

export default function SearchUserPage() {
    const [page, setPage] = useState(0);
    const [order, setOrder] = useState('desc');
    const [orderBy, setOrderBy] = useState('links');
    const [filterName, setFilterName] = useState('');
    const [rowsPerPage, setRowsPerPage] = useState(25);

    const [errorToast, setErrorToast] = useState({ open: false, message: '' });
    const [successToast, setSuccessToast] = useState({ open: false, message: '', url: '' });

    const [issues, setIssues] = useState([]);
    const [users, setUsers] = useState([]);
    const [username, setUsername] = useState('');
    const [selectedUser, setSelectedUser] = useState(null);

    // State mangagement for user flow
    const [isSearchingForUser, setIsSearchingForUser] = useState(false);
    const [isSavingRepos, setIsSavingRepos] = useState(false);
    const [showUsersGrid, setShowUsersGrid] = useState(false);
    const [showRepoIssues, setShowRepoIssues] = useState(false);

    useEffect(() => {
        // On componenet mount, reload user results data, if there
        try {
            loadResults();
        } catch (e) {
            // Catch any unexpected errors (as we cant control whats in local storage)
            console.error('Error loading user results', e);
            setErrorToast({ open: true, message: e.toString() });
        }
    }, []);

    const handleSort = (event, id) => {
        const isAsc = orderBy === id && order === 'asc';
        if (id !== '') {
            setOrder(isAsc ? 'desc' : 'asc');
            setOrderBy(id);
        }
    };

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        setPage(0);
        setRowsPerPage(parseInt(event.target.value, 10));
    };

    const handleFilterByName = (event) => {
        setPage(0);
        setFilterName(event.target.value);
    };

    const issuesFiltered = applyFilter({
        inputData: issues,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const handleCloseErrorToast = () => {
        setErrorToast({ open: false, message: '' });
    };

    const handleCloseSuccessToast = () => {
        setSuccessToast({ open: false, message: '' });
    };

    const handleSearch = async () => {
        if (isSearchingForUser || isSavingRepos) {
            return;
        }

        setIsSearchingForUser(true);
        setShowUsersGrid(false);
        setShowRepoIssues(false);
        console.log(`Starting github search for ${username}`);

        try {
            const res = await searchGithubUsers(username);
            if (res === null) {
                setUsers([]);
            } else {
                setUsers(res);
            }
        } catch (e) {
            console.error('Error searching for users', e);
            setErrorToast({ open: true, message: e.toString() });
            setIsSearchingForUser(false);
            setShowUsersGrid(false);
            return;
        }

        // Search has completed
        setIsSearchingForUser(false);
        setShowUsersGrid(true);
    };

    const handleUserSelect = async (login) => {
        setSelectedUser(login);
        console.log(`Selected user: ${login}..getting repos`);
        setIsSavingRepos(true);

        clearResults();

        //  Find Repos for user
        var repositories;
        try {
            const res = await searchGithubUsersRepo(login);
            repositories = res;
        } catch (e) {
            console.error('Error searching for repos', e);
            setErrorToast({ open: true, message: e.toString() });

            setShowUsersGrid(true);
            setIsSavingRepos(false);
            setShowRepoIssues(false);
            return;
        }

        // Save repos to DB
        try {
            await saveRepos(repositories);
            console.log(repositories.length, 'repos saved to db');
        } catch (e) {
            console.log('Error saving repos to db', e);
            setErrorToast({ open: true, message: e.toString() });

            setShowUsersGrid(true);
            setIsSavingRepos(false);
            setShowRepoIssues(false);
            return;
        }

        // Pull states of repos
        try {
            const res = await searchRepoLinksForUser(login);
            if (res === null) {
                setIssues([]);
            } else {
                setIssues(res);
            }

            console.log('Search completed!');
        } catch (e) {
            console.error('Error searching for repo links', e);
            setErrorToast({ open: true, message: e.toString() });

            setShowUsersGrid(true);
            setIsSavingRepos(false);
            setShowRepoIssues(false);
            return;
        }

        // Update stats
        setShowUsersGrid(false);
        setIsSavingRepos(false);
        setShowRepoIssues(true);

        setSuccessToast({
            open: true,
            message: `Found ${repositories.length} repo(s). Please refresh after a few seconds to see results.`,
        });

        // At this point we should snapshot data and save to localStorage for an easy lookup/refresh
        saveResults(login);
    };

    const saveResults = (user) => {
        localStorage.setItem('selectedUser', user);
    };

    const clearResults = () => {
        localStorage.removeItem('selectedUser');
        setSelectedUser(null);
    };

    const loadResults = async () => {
        const selectedUser = localStorage.getItem('selectedUser');
        if (selectedUser) {
            setSelectedUser(selectedUser);
            try {
                const res = await searchRepoLinksForUser(selectedUser);
                if (res === null) {
                    setIssues([]);
                } else {
                    setIssues(res);
                }

                console.log('Search completed!');
            } catch (e) {
                console.error('Error searching for repo links', e);
                setErrorToast({ open: true, message: e.toString() });

                setShowUsersGrid(true);
                setIsSavingRepos(false);
                setShowRepoIssues(false);
                return;
            }

            setShowUsersGrid(false);
            setIsSavingRepos(false);
            setShowRepoIssues(true);
        }
    };

    const refreshIssues = async () => {
        if (isSavingRepos && isSearchingForUser) {
            return;
        }

        console.log('Refreshing issues');
        // console.log('Selected user:', selectedUser);
        const selectedUser = localStorage.getItem('selectedUser');

        try {
            const res = await searchRepoLinksForUser(selectedUser);
            if (res === null) {
                setIssues([]);
            } else {
                setIssues(res);
            }

            console.log('Search completed!');
        } catch (e) {
            console.error('Error searching for repo links', e);
            setErrorToast({ open: true, message: e.toString() });

            setShowUsersGrid(true);
            setIsSavingRepos(false);
            setShowRepoIssues(false);
            return;
        }

        setShowRepoIssues(true);
        setSuccessToast({
            open: true,
            message: `Refreshed ${issues.length} repo(s).`,
        });
    };

    const handleCreateIssue = async (repoId) => {
        console.log('Creating issue for repo:', repoId);
        try {
            const res = await createIssue(repoId);
            setSuccessToast({
                open: true,
                message: `Issue created:`,
                url: res.url,
            });
        } catch (e) {
            console.error('Error creating issue', e);
            setErrorToast({ open: true, message: e.toString() });
        }
    };

    const notFound = !issuesFiltered.length && !!filterName;

    return (
        <Container>
            <Stack direction="column" alignItems="left" justifyContent="space-between" mb={5}>
                <Stack mb={1}>
                    <Typography variant="h4">Search for new Repositories</Typography>
                </Stack>
                <Typography variant="p">
                    Here you can search GitHub for new repositories to add to the GooGL-Bye queue.
                </Typography>
            </Stack>

            <SearchUserForm
                setUsername={setUsername}
                handleSearch={handleSearch}
                username={username}
                isSearching={isSearchingForUser}
            />

            {showUsersGrid && users.length > 0 ? (
                <UserGrid
                    isFindingRepos={isSavingRepos}
                    users={users}
                    onUserSelect={handleUserSelect}
                />
            ) : showUsersGrid ? (
                <Card sx={{ mt: 4, p: 3, textAlign: 'center' }}>
                    <Typography variant="h6" color="text.secondary">
                        No users found
                    </Typography>
                    <Typography variant="body2" sx={{ mt: 1 }}>
                        We couldn't find any GitHub users matching your search. Please try a
                        different username.
                    </Typography>
                </Card>
            ) : null}

            {showRepoIssues && (
                <Card>
                    <IssueTableToolbar
                        numSelected={0}
                        filterName={filterName}
                        onFilterName={handleFilterByName}
                        refreshIssues={refreshIssues}
                    />

                    <Scrollbar>
                        <TableContainer sx={{ overflow: 'unset' }}>
                            <Table sx={{ minWidth: 800 }}>
                                <IssueTableHead
                                    order={order}
                                    orderBy={orderBy}
                                    rowCount={issues.length}
                                    numSelected={0}
                                    onRequestSort={handleSort}
                                    headLabel={[
                                        { id: 'name', label: 'Repository Name' },
                                        { id: 'author', label: 'Author' },
                                        { id: 'language', label: 'Language' },
                                        { id: 'stars', label: 'Stars', align: 'center' },
                                        { id: 'forks', label: 'Forks', align: 'center' },
                                        { id: 'state', label: 'Status' },
                                        { id: 'links', label: 'Links' },
                                        { id: '' },
                                    ]}
                                />
                                <TableBody>
                                    {issues.length > 0 ? (
                                        <>
                                            {issuesFiltered
                                                .slice(
                                                    page * rowsPerPage,
                                                    page * rowsPerPage + rowsPerPage
                                                )
                                                .map((row) => (
                                                    <IssueTableRow
                                                        key={row.id}
                                                        name={row.name}
                                                        author={row.author}
                                                        language={row.language}
                                                        stars={row.stars}
                                                        forks={row.forks}
                                                        ghUrl={row.ghUrl}
                                                        state={row.state}
                                                        avatarUrl={row.avatarUrl}
                                                        lastCommit={row.lastCommit}
                                                        issues={row.links}
                                                        errorMsg={row.errorMsg}
                                                        selected={0}
                                                        handleCreateIssue={() =>
                                                            handleCreateIssue(row.id)
                                                        }
                                                    />
                                                ))}
                                            <TableEmptyRows
                                                height={77}
                                                emptyRows={emptyRows(
                                                    page,
                                                    rowsPerPage,
                                                    issues.length
                                                )}
                                            />
                                            {notFound && <TableNoData query={filterName} />}
                                        </>
                                    ) : (
                                        <TableRow>
                                            <TableCell colSpan={8} align="center">
                                                <Typography variant="subtitle1">
                                                    No repos found for user {selectedUser}.
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Scrollbar>
                    <TablePagination
                        page={page}
                        component="div"
                        count={issues.length}
                        rowsPerPage={rowsPerPage}
                        onPageChange={handleChangePage}
                        rowsPerPageOptions={[5, 10, 25, 50]}
                        onRowsPerPageChange={handleChangeRowsPerPage}
                    />
                </Card>
            )}

            <ErrorToast
                open={errorToast.open}
                message={errorToast.message}
                onClose={handleCloseErrorToast}
            />
            <SuccessToast
                open={successToast.open}
                message={successToast.message}
                onClose={handleCloseSuccessToast}
                url={successToast.url}
            />
        </Container>
    );
}

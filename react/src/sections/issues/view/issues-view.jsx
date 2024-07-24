import { useState, useEffect } from 'react';

import {
    Stack,
    Table,
    Button,
    Container,
    TableBody,
    TableRow,
    TableCell,
    TableContainer,
    TablePagination,
    Typography,
    Card,
    Backdrop,
    CircularProgress,
} from '@mui/material';

import { generateMockReposLinks } from 'src/_mock/repoLinks';

import Scrollbar from 'src/components/scrollbar';

import ErrorToast from 'src/components/toast/errorToast';

import TableNoData from '../table-no-data';
import RepoTableRow from '../issues-table-row';
import RepoTableHead from '../issues-table-head';
import TableEmptyRows from '../table-empty-rows';
import RepoTableToolbar from '../issues-table-toolbar';
import { emptyRows, applyFilter, getComparator } from '../utils';

import { searchRepoLinks } from 'src/api/client';

// ----------------------------------------------------------------------

export default function IssuesPage() {
    const [page, setPage] = useState(0);
    const [order, setOrder] = useState('asc');
    const [selected, setSelected] = useState([]);
    const [orderBy, setOrderBy] = useState('name');
    const [filterName, setFilterName] = useState('');
    const [rowsPerPage, setRowsPerPage] = useState(10);

    const [errorToast, setErrorToast] = useState({ open: false, message: '' });
    const [repos, setRepos] = useState([]);

    const [isLoading, setIsLoading] = useState([]);

    useEffect(() => {
        async function fetchData() {
            await loadIssues();
        }
        fetchData();
    }, []);

    const refreshIssues = async () => {
        await loadIssues();
    };

    async function loadIssues() {
        console.log('Getting repos state from db');
        setIsLoading(true);

        try {
            const res = await searchRepoLinks();
            if (res === null) {
                setRepos([]);
            } else {
                setRepos(res);
            }
            console.log('Search completed!');
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
        }
        setSelected([]);
        setIsLoading(false);

        // setTimeout(() => {
        //     setRepos(generateMockReposLinks);
        //     setSelected([]);
        //     setIsLoading(false);
        //     console.log('Lookup completed!');
        // }, 500);
    }

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
        Bad;
    };

    const dataFiltered = applyFilter({
        inputData: repos,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const handleCloseErrorToast = () => {
        setErrorToast({ open: false, message: '' });
    };

    const notFound = !dataFiltered.length && !!filterName;

    return (
        <Container>
            <Stack direction="column" alignItems="left" justifyContent="space-between" mb={5}>
                <Stack direction="row" alignItems="center" justifyContent="space-between" mb={1}>
                    <Typography variant="h4">Processed Repositories</Typography>
                    <Button onClick={refreshIssues} variant="contained" color="primary">
                        Refresh
                    </Button>
                </Stack>
                <Typography variant="p">
                    Here you can find repositories in the system, which have goo.gl links still in their code
                    base. You can perform actions such as raising automated issues or PRs from this page.
                    (Coming soon)
                </Typography>
            </Stack>

            <Card>
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
                    open={isLoading}
                >
                    <CircularProgress color="inherit" />
                </Backdrop>
                <RepoTableToolbar
                    numSelected={selected.length}
                    filterName={filterName}
                    onFilterName={handleFilterByName}
                />

                <Scrollbar>
                    <TableContainer sx={{ overflow: 'unset' }}>
                        <Table sx={{ minWidth: 800 }}>
                            <RepoTableHead
                                order={order}
                                orderBy={orderBy}
                                rowCount={repos.length}
                                numSelected={selected.length}
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
                                {repos.length > 0 ? (
                                    <>
                                        {dataFiltered
                                            .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                            .map((row) => (
                                                <RepoTableRow
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
                                                    selected={selected.indexOf(row.name) !== -1}
                                                />
                                            ))}
                                        <TableEmptyRows
                                            height={77}
                                            emptyRows={emptyRows(page, rowsPerPage, repos.length)}
                                        />
                                        {notFound && <TableNoData query={filterName} />}
                                    </>
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={7} align="center">
                                            <Typography variant="subtitle1">
                                                No repositories found. Use the search to find repositories.
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
                    count={repos.length}
                    rowsPerPage={rowsPerPage}
                    onPageChange={handleChangePage}
                    rowsPerPageOptions={[5, 10, 25]}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                />
            </Card>
            <ErrorToast
                open={errorToast.open}
                message={errorToast.message}
                onClose={handleCloseErrorToast}
            />
        </Container>
    );
}

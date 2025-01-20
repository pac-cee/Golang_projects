import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  IconButton,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Typography,
  Chip,
  TextField,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { Category } from '../types';
import { useAppContext } from '../context/AppContext';
import { categoryApi } from '../services/api';
import CategoryForm from '../components/forms/CategoryForm';

const Categories: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | undefined>();
  const [search, setSearch] = useState('');

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        const categories = await categoryApi.getAll();
        dispatch({ type: 'SET_CATEGORIES', payload: categories });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error fetching categories' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    };

    fetchCategories();
  }, [dispatch]);

  const handleAddCategory = async (
    category: Omit<Category, 'id' | 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const newCategory = await categoryApi.create(category);
      dispatch({ type: 'ADD_CATEGORY', payload: newCategory });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error adding category' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleUpdateCategory = async (
    category: Omit<Category, 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const updatedCategory = await categoryApi.update(category.id, category);
      dispatch({ type: 'UPDATE_CATEGORY', payload: updatedCategory });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error updating category' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleDeleteCategory = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this category?')) {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        await categoryApi.delete(id);
        dispatch({ type: 'DELETE_CATEGORY', payload: id });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error deleting category' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    }
  };

  const filteredCategories = state.categories.filter((category) =>
    category.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <Box>
      <Grid container spacing={3}>
        <Grid item xs={12} display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h4">Categories</Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => {
              setSelectedCategory(undefined);
              setOpenForm(true);
            }}
          >
            Add Category
          </Button>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <TextField
                fullWidth
                size="small"
                label="Search Categories"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <List>
            {filteredCategories.map((category) => (
              <ListItem
                key={category.id}
                component={Card}
                sx={{ mb: 2, display: 'block' }}
              >
                <CardContent>
                  <Grid container alignItems="center" spacing={2}>
                    <Grid item xs>
                      <Typography variant="h6">{category.name}</Typography>
                      <Box sx={{ mt: 1 }}>
                        {category.subcategories.map((subcategory, index) => (
                          <Chip
                            key={index}
                            label={subcategory}
                            size="small"
                            variant="outlined"
                            sx={{ mr: 1, mb: 1 }}
                          />
                        ))}
                      </Box>
                    </Grid>
                    <Grid item>
                      <ListItemSecondaryAction>
                        <IconButton
                          edge="end"
                          onClick={() => {
                            setSelectedCategory(category);
                            setOpenForm(true);
                          }}
                        >
                          <EditIcon />
                        </IconButton>
                        <IconButton
                          edge="end"
                          onClick={() => handleDeleteCategory(category.id)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </Grid>
                  </Grid>
                </CardContent>
              </ListItem>
            ))}
          </List>
        </Grid>
      </Grid>

      <CategoryForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedCategory(undefined);
        }}
        category={selectedCategory}
        onSubmit={selectedCategory ? handleUpdateCategory : handleAddCategory}
      />
    </Box>
  );
};

export default Categories;
